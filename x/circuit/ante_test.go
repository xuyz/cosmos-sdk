package circuit

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	params "github.com/cosmos/cosmos-sdk/x/params"
)

type tx []sdk.Msg

var _ sdk.Tx = tx{}

func (tx tx) ValidateBasic() sdk.Error { return nil }

func (tx tx) GetMsgs() []sdk.Msg { return tx }

type msg struct{}

func (msg) ValidateBasic() sdk.Error { return nil }

func (msg) GetSignBytes() []byte { return nil }

func (msg) GetSigners() []sdk.AccAddress { return nil }

type msg1 struct{ msg }

var _ sdk.Msg = msg1{}

func (msg1) Route() string { return "msg1" }

func (msg1) Type() string { return "msg1" }

type msg2 struct{ msg }

var _ sdk.Msg = msg2{}

func (msg2) Route() string { return "msg2" }

func (msg2) Type() string { return "msg2" }

func testMsg(t *testing.T, ctx sdk.Context, k Keeper, msg sdk.Msg, initial bool) {
	ante := NewAnteHandler(k)

	_, _, abort := ante(ctx, tx{msg}, false)
	require.Equal(t, initial, abort)

	require.NotPanics(t, func() { k.space.SetWithSubkey(ctx, MsgRouteKey, []byte(msg.Type()), true) }, "panic setting breaker")
	_, _, abort = ante(ctx, tx{msg}, false)
	require.Equal(t, true, abort)

	require.NotPanics(t, func() { k.space.SetWithSubkey(ctx, MsgRouteKey, []byte(msg.Type()), false) }, "panic setting breaker")
	_, _, abort = ante(ctx, tx{msg}, false)
	require.Equal(t, false, abort)
}

func TestAnteHandler(t *testing.T) {
	ctx, space, _ := params.DefaultTestComponents(t)

	k := NewKeeper(space)

	data := GenesisState{
		MsgRoutes: []string{"msg2"},
	}

	InitGenesis(ctx, k, data)

	testMsg(t, ctx, k, msg1{}, false)
	testMsg(t, ctx, k, msg2{}, true)
}

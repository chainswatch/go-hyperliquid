package hyperliquid

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func ToTypedSig(r [32]byte, s [32]byte, v byte) RsvSignature {
	return RsvSignature{
		R: hexutil.Encode(r[:]),
		S: hexutil.Encode(s[:]),
		V: v,
	}
}

func ArrayAppend(data []byte, toAppend []byte) []byte {
	return append(data, toAppend...)
}

func HexToBytes(addr string) []byte {
	if strings.HasPrefix(addr, "0x") {
		fAddr := strings.Replace(addr, "0x", "", 1)
		b, _ := hex.DecodeString(fAddr)
		return b
	} else {
		b, _ := hex.DecodeString(addr)
		return b
	}
}

func OrderWiresToOrderAction(orders []OrderWire, grouping Grouping) PlaceOrderAction {
	return PlaceOrderAction{
		Type:     "order",
		Grouping: grouping,
		Orders:   orders,
	}
}

func OrderRequestToWire(req OrderRequest, meta map[string]AssetInfo, isSpot bool) OrderWire {
	info := meta[req.Coin]
	var assetId, maxDecimals int
	if isSpot {
		// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/asset-ids
		assetId = info.AssetId + 10000
		maxDecimals = SPOT_MAX_DECIMALS
	} else {
		assetId = info.AssetId
		maxDecimals = PERP_MAX_DECIMALS
	}
	return OrderWire{
		Asset:      assetId,
		IsBuy:      req.IsBuy,
		LimitPx:    PriceToWire(req.LimitPx, maxDecimals, info.SzDecimals),
		SizePx:     SizeToWire(req.Sz, info.SzDecimals),
		ReduceOnly: req.ReduceOnly,
		OrderType:  req.OrderType,
		Cloid:      req.Cloid,
	}
}

// Format the float with custom decimal places, default is 6 (perp), 8 (spot).
// https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api/tick-and-lot-size
func FloatToWire(x float64, maxDecimals int, szDecimals int) string {
	bigf := big.NewFloat(x)
	var maxDecSz uint
	intPart, _ := bigf.Int64()
	intSize := len(strconv.FormatInt(intPart, 10))
	if intSize >= maxDecimals {
		maxDecSz = 0
	} else {
		maxDecSz = uint(maxDecimals - intSize)
	}
	x, _ = bigf.Float64()
	rounded := fmt.Sprintf("%.*f", maxDecSz, x)
	if strings.Contains(rounded, ".") {
		for strings.HasSuffix(rounded, "0") {
			rounded = strings.TrimSuffix(rounded, "0")
		}
	}
	if strings.HasSuffix(rounded, ".") {
		rounded = strings.TrimSuffix(rounded, ".")
	}
	return rounded
}

// fastPow10 returns 10^exp as a float64. For our purposes exp is small.
func pow10(exp int) float64 {
	var res float64 = 1
	for i := 0; i < exp; i++ {
		res *= 10
	}
	return res
}

// PriceToWire converts a price value to its string representation per Hyperliquid rules.
// It enforces:
//   - At most 5 significant figures,
//   - And no more than (maxDecimals - szDecimals) decimal places.
//
// Integer prices are returned as is.
func PriceToWire(x float64, _ int, _ int) string {
	// The library previously tried to enforce tick/lot size restrictions.
	// This behaviour was opinionated and has been removed.  Now we simply
	// return the string representation of the provided value without any
	// rounding or truncation.
	return strconv.FormatFloat(x, 'f', -1, 64)
}

// SizeToWire converts a size value to its string representation,
// rounding it to exactly szDecimals decimals.
// Integer sizes are returned without decimals.
func SizeToWire(x float64, _ int) string {
	// As with PriceToWire, no automatic rounding is applied.  The caller is
	// responsible for providing values that already respect the desired
	// precision.
	return strconv.FormatFloat(x, 'f', -1, 64)
}

// To sign raw messages via EIP-712
func StructToMap(strct any) (res map[string]interface{}, err error) {
	a, err := json.Marshal(strct)
	if err != nil {
		return map[string]interface{}{}, err
	}
	json.Unmarshal(a, &res)
	return res, nil
}

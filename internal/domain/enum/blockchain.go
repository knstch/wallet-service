package enum

import (
	public "github.com/knstch/wallets-api/public"

	"github.com/knstch/subtrack-libs/enum"
)

func ConvertPublicNetworkToService(network public.Network) enum.Network {
	switch network {
	case public.Network_POLYGON:
		return enum.PolygonNetwork
	case public.Network_BSC:
		return enum.BscNetwork
	default:
		return enum.UnknownNetwork
	}
}

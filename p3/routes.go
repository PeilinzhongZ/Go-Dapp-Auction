package p3

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
	Route{
		"Upload",
		"GET",
		"/upload",
		Upload,
	},
	Route{
		"UploadBlock",
		"GET",
		"/block/{height}/{hash}",
		UploadBlock,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
	Route{
		"UploadPeerMap",
		"GET",
		"/peerMap",
		UploadPeerMap,
	},
	Route{
		"Canonical",
		"GET",
		"/canonical",
		Canonical,
	},
	Route{
		"Post",
		"POST",
		"/postItem",
		PostItem,
	},
	Route{
		"ListItem",
		"GET",
		"/listItem",
		ListItem,
	},
	Route{
		"Post",
		"POST",
		"/postBid",
		PostBid,
	},
	Route{
		"Finalize",
		"GET",
		"/finalize",
		FinalizeAuction,
	},
	Route{
		"Validate",
		"GET",
		"/check",
		CheckResult,
	},
}

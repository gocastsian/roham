package roham.rego

import rego.v1

role_admin := uint8(1)

default rule_check_request_only := false

rule_check_request_only if {
    some request in ["GetCapabilities", "GetMap", "GetFeatureInfo", "DescribeLayer", "GetLegendGraphic"]
    input.request == request
    some role in input.role
    role == role_admin
}

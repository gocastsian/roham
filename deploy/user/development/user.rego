package user.authz

default allow = false

# Define role constants
role_admin := 1

# Admin can access everything
allow if {
    input.user.role == role_admin
}

# Regular users can access if the request has query parameter "request=GetMap"
allow if {
    input.request.query.request == "GetMap"
}
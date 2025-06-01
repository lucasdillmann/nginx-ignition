import { UserAccessLevel } from "./UserAccessLevel"

export default interface UserPermissions {
    hosts: UserAccessLevel
    streams: UserAccessLevel
    certificates: UserAccessLevel
    logs: UserAccessLevel
    integrations: UserAccessLevel
    accessLists: UserAccessLevel
    settings: UserAccessLevel
    users: UserAccessLevel
    nginxServer: UserAccessLevel
}

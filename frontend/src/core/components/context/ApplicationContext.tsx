import React from "react";
import UserService from "../../../domain/user/UserService";
import UserResponse from "../../../domain/user/model/UserResponse";

export interface ApplicationContextData {
    user?: UserResponse,
}

export async function startApplicationContext(): Promise<ApplicationContextData> {
    const user = await new UserService().current()

    return {
        user,
    }
}

const ApplicationContext = React.createContext<ApplicationContextData>({})
export default ApplicationContext

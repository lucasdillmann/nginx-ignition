import React from "react";
import UserService from "../../../domain/user/UserService";
import UserResponse from "../../../domain/user/model/UserResponse";

export interface AppContextData {
    user?: UserResponse,
}

export async function loadAppContextData(): Promise<AppContextData> {
    const user = await new UserService().current()

    return {
        user,
    }
}

const AppContext = React.createContext<AppContextData>({})
export default AppContext

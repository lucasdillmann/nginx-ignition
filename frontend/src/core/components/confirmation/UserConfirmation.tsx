import {Modal} from "antd";
import React from "react";

class UserConfirmation {
    ask(message: React.ReactNode): Promise<void> {
        const messageContainer = (
            <div style={{ margin: "0 0 15px" }}>
                {message}
            </div>
        )

        return new Promise((resolve) => {
            Modal.confirm({
                title: 'Beware, young padawan',
                content: messageContainer,
                cancelText: "No",
                okText: "Yes",
                okButtonProps: {
                    color: "danger",
                    variant: "solid",
                },
                onOk: () => resolve(),
            });
        });
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new UserConfirmation()

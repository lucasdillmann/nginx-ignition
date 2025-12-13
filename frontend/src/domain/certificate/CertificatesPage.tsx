import React from "react"
import AppShellContext from "../../core/components/shell/AppShellContext"

export default class CertificatesPage extends React.Component {
    constructor(props: any) {
        super(props)
    }

    componentDidMount() {
        AppShellContext.get().updateConfig({
            title: "Certificates",
            subtitle: "",
            actions: [],
        })
    }

    render() {
        return <></>
    }
}

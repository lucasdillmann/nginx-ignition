import React from "react"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import CacheService from "./CacheService"
import { Form, FormInstance, Tabs } from "antd"
import Preloader from "../../core/components/preloader/Preloader"
import FormLayout from "../../core/components/form/FormLayout"
import ValidationResult from "../../core/validation/ValidationResult"
import ModalPreloader from "../../core/components/preloader/ModalPreloader"
import Notification from "../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import AppShellContext, { ShellAction } from "../../core/components/shell/AppShellContext"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import DeleteCacheAction from "./actions/DeleteCacheAction"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import CacheRequest from "./model/CacheRequest"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import AccessControl from "../../core/components/accesscontrol/AccessControl"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { cacheFormDefaults } from "./CacheFormDefaults"
import GeneralTab from "./tabs/GeneralTab"
import AdvancedTab from "./tabs/AdvancedTab"
import "./CacheFormPage.css"

interface CacheFormState {
    formValues: CacheRequest
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
    error?: Error
}

export default class CacheFormPage extends React.Component<unknown, CacheFormState> {
    private readonly service: CacheService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>
    private cacheId?: string

    constructor(props: any) {
        super(props)
        const cacheId = routeParams().id
        this.formRef = React.createRef()
        this.cacheId = cacheId === "new" ? undefined : cacheId
        this.service = new CacheService()
        this.saveModal = new ModalPreloader()
        this.state = {
            formValues: cacheFormDefaults(),
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
        }
    }

    private submit() {
        const { formValues } = this.state
        this.saveModal.show("Hang on tight", "We're saving the cache configuration")
        this.setState({ validationResult: new ValidationResult() })

        const action =
            this.cacheId === undefined
                ? this.service.create(formValues).then(response => this.updateId(response.id))
                : this.service.updateById(this.cacheId, formValues)

        action.then(() => this.handleSuccess()).catch(error => this.handleError(error))
    }

    private updateId(id: string) {
        this.cacheId = id
        navigateTo(`/caches/${id}`, true)
        this.updateShellConfig(true)
    }

    private handleSuccess() {
        this.saveModal.close()
        Notification.success("Cache configuration saved", "The cache configuration was saved successfully")
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        this.saveModal.close()
        Notification.error("That didn't work", "Please check the form to see if everything seems correct")
    }

    private handleChange(newValues: CacheRequest) {
        this.setState({
            formValues: {
                ...newValues,
            },
        })
    }

    private buildTabs(): any[] {
        const { validationResult } = this.state

        return [
            {
                key: "general",
                label: "General",
                forceRender: true,
                children: <GeneralTab validationResult={validationResult} />,
            },
            {
                key: "advanced",
                label: "Advanced",
                forceRender: true,
                children: <AdvancedTab validationResult={validationResult} />,
            },
        ]
    }

    private renderForm() {
        const { formValues } = this.state

        return (
            <Form<CacheRequest>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                onValuesChange={(_, formValues) => this.handleChange(formValues)}
                initialValues={formValues}
            >
                <Tabs items={this.buildTabs()} destroyOnHidden={false} />
            </Form>
        )
    }

    private async delete() {
        if (this.cacheId === undefined) return

        return DeleteCacheAction.execute(this.cacheId).then(() => navigateTo("/caches"))
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.caches)) {
            enableActions = false
        }

        const actions: ShellAction[] = [
            {
                description: "Save",
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.cacheId !== undefined)
            actions.unshift({
                description: "Delete",
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: "Cache configuration details",
            subtitle: "Full details of the nginx content cache configuration",
            actions,
        })
    }

    componentDidMount() {
        if (this.cacheId === undefined) {
            this.setState({ loading: false })
            this.updateShellConfig(true)
            return
        }

        this.service
            .getById(this.cacheId!!)
            .then(cache => {
                if (cache === undefined) this.setState({ loading: false, notFound: true })
                else {
                    this.setState({ loading: false, formValues: cache })
                    this.updateShellConfig(true)
                }
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })

        this.updateShellConfig(false)
    }

    render() {
        const { loading, notFound, error } = this.state

        if (error !== undefined) return EmptyStates.FailedToFetch
        if (notFound) return EmptyStates.NotFound
        if (loading) return <Preloader loading />

        return (
            <AccessControl
                requiredAccessLevel={UserAccessLevel.READ_ONLY}
                permissionResolver={permissions => permissions.caches}
            >
                {this.renderForm()}
            </AccessControl>
        )
    }
}

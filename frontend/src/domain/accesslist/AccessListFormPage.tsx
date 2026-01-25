import React from "react"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import AccessListService from "./AccessListService"
import { Form, FormInstance, Input, Select, Switch } from "antd"
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
import DeleteAccessListAction from "./actions/DeleteAccessListAction"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import AccessListFormValues from "./model/AccessListFormValues"
import AccessListConverter from "./AccessListConverter"
import AccessListEntrySets from "./components/AccessListEntrySets"
import AccessListCredentials from "./components/AccessListCredentials"
import { AccessListOutcome } from "./model/AccessListRequest"
import "./AccessListFormPage.css"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import AccessControl from "../../core/components/accesscontrol/AccessControl"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { accessListFormDefaults } from "./AccessListFormDefaults"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../core/i18n/I18n"

interface AccessListFormState {
    formValues: AccessListFormValues
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
    error?: Error
}

export default class AccessListFormPage extends React.Component<unknown, AccessListFormState> {
    private readonly service: AccessListService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>
    private accessListId?: string

    constructor(props: any) {
        super(props)
        const accessListId = routeParams().id
        this.formRef = React.createRef()
        this.accessListId = accessListId === "new" ? undefined : accessListId
        this.service = new AccessListService()
        this.saveModal = new ModalPreloader()
        this.state = {
            formValues: accessListFormDefaults(),
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
        }
    }

    private submit() {
        const { formValues } = this.state
        this.saveModal.show(MessageKey.CommonHangOnTight, {
            id: MessageKey.CommonSavingType,
            params: { type: MessageKey.CommonAccessList },
        })
        this.setState({ validationResult: new ValidationResult() })

        const data = AccessListConverter.toRequest(formValues)
        const action =
            this.accessListId === undefined
                ? this.service.create(data).then(response => this.updateId(response.id))
                : this.service.updateById(this.accessListId, data)

        action.then(() => this.handleSuccess()).catch(error => this.handleError(error))
    }

    private updateId(id: string) {
        this.accessListId = id
        navigateTo(`/access-lists/${id}`, true)
        this.updateShellConfig(true)
    }

    private handleSuccess() {
        this.saveModal.close()
        Notification.success(
            { id: MessageKey.CommonTypeSaved, params: { type: MessageKey.CommonAccessList } },
            MessageKey.CommonSuccessMessage,
        )
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        this.saveModal.close()
        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private removeEntry(index: number) {
        const { formValues } = this.state
        const { entries } = formValues

        let priority = 0
        const updatedValues = entries
            .filter((_, itemIndex) => itemIndex !== index)
            .map(route => ({
                ...route,
                priority: priority++,
            }))

        this.formRef.current?.setFieldValue("entries", updatedValues)
        this.setState({
            formValues: {
                ...formValues,
                entries: updatedValues,
            },
        })
    }

    private handleChange(newValues: AccessListFormValues) {
        const entries = newValues.entries.sort((left, right) => (left.priority > right.priority ? 1 : -1))

        this.setState({
            formValues: {
                ...newValues,
                entries,
            },
        })
    }

    private renderForm() {
        const { formValues, validationResult } = this.state

        return (
            <Form<AccessListFormValues>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                onValuesChange={(_, formValues) => this.handleChange(formValues)}
                initialValues={formValues}
            >
                <h2 className="access-lists-form-section-name">
                    <I18n id={MessageKey.CommonGeneral} />
                </h2>
                <p className="access-lists-form-section-help-text">
                    <I18n id={MessageKey.FrontendAccesslistSectionGeneralDescription} />
                </p>
                <Form.Item
                    className="access-lists-form-name"
                    name="name"
                    validateStatus={validationResult.getStatus("name")}
                    help={validationResult.getMessage("name")}
                    label={<I18n id={MessageKey.CommonName} />}
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    className="access-lists-form-realm-name"
                    name="realm"
                    validateStatus={validationResult.getStatus("realm")}
                    help={validationResult.getMessage("realm")}
                    label={<I18n id={MessageKey.FrontendAccesslistRealmName} />}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    className="access-lists-form-satisfy-all"
                    name="satisfyAll"
                    validateStatus={validationResult.getStatus("satisfyAll")}
                    help={validationResult.getMessage("satisfyAll")}
                    label={<I18n id={MessageKey.CommonMode} />}
                    required
                >
                    <Select>
                        <Select.Option value={true}>
                            <I18n id={MessageKey.FrontendAccesslistModeSatisfyAll} />
                        </Select.Option>
                        <Select.Option value={false}>
                            <I18n id={MessageKey.FrontendAccesslistModeSatisfyAny} />
                        </Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    className="access-lists-form-default-outcome"
                    name="defaultOutcome"
                    validateStatus={validationResult.getStatus("defaultOutcome")}
                    help={validationResult.getMessage("defaultOutcome")}
                    label={<I18n id={MessageKey.FrontendAccesslistDefaultOutcome} />}
                    required
                >
                    <Select>
                        <Select.Option value={AccessListOutcome.ALLOW}>
                            <I18n id={MessageKey.FrontendAccesslistOutcomeAllow} />
                        </Select.Option>
                        <Select.Option value={AccessListOutcome.DENY}>
                            <I18n id={MessageKey.FrontendAccesslistOutcomeDeny} />
                        </Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    className="access-lists-form-forward-headers"
                    name="forwardAuthenticationHeader"
                    validateStatus={validationResult.getStatus("forwardAuthenticationHeader")}
                    help={validationResult.getMessage("forwardAuthenticationHeader")}
                    label={<I18n id={MessageKey.FrontendAccesslistForwardAuthHeaders} />}
                    required
                >
                    <Switch />
                </Form.Item>

                <h2 className="access-lists-form-section-name">
                    <I18n id={MessageKey.FrontendAccesslistSectionCredentials} />
                </h2>
                <p className="access-lists-form-section-help-text">
                    <I18n id={MessageKey.FrontendAccesslistSectionCredentialsDescription} />
                </p>
                <AccessListCredentials validationResult={validationResult} />

                <h2 className="access-lists-form-section-name">
                    <I18n id={MessageKey.FrontendAccesslistSectionEntrySets} />
                </h2>
                <p className="access-lists-form-section-help-text">
                    <I18n id={MessageKey.FrontendAccesslistSectionEntrySetsDescription} />
                </p>
                <AccessListEntrySets
                    entrySets={formValues.entries}
                    validationResult={validationResult}
                    onRemove={index => this.removeEntry(index)}
                />
            </Form>
        )
    }

    private async delete() {
        if (this.accessListId === undefined) return

        return DeleteAccessListAction.execute(this.accessListId).then(() => navigateTo("/access-lists"))
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.accessLists)) {
            enableActions = false
        }

        const actions: ShellAction[] = [
            {
                description: MessageKey.CommonSave,
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.accessListId !== undefined)
            actions.unshift({
                description: MessageKey.CommonDelete,
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: MessageKey.FrontendAccesslistFormTitle,
            subtitle: MessageKey.FrontendAccesslistFormSubtitle,
            actions,
        })
    }

    componentDidMount() {
        if (this.accessListId === undefined) {
            this.setState({ loading: false })
            this.updateShellConfig(true)
            return
        }

        this.service
            .getById(this.accessListId!!)
            .then(accessListDetails => {
                if (accessListDetails === undefined) this.setState({ loading: false, notFound: true })
                else {
                    this.setState({ loading: false, formValues: AccessListConverter.toFormValues(accessListDetails) })
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
                permissionResolver={permissions => permissions.accessLists}
            >
                {this.renderForm()}
            </AccessControl>
        )
    }
}

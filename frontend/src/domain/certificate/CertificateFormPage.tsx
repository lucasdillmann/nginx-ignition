import React from "react";
import ShellAwareComponent, {ShellConfig} from "../../core/components/shell/ShellAwareComponent";
import ValidationResult from "../../core/validation/ValidationResult";
import CertificateService from "./CertificateService";
import AvailableProviderResponse from "./model/AvailableProviderResponse";
import Preloader from "../../core/components/preloader/Preloader";
import {Form, FormInstance, Select} from "antd";
import If from "../../core/components/flowcontrol/If";
import FormLayout from "../../core/components/form/FormLayout";
import DynamicInput from "./components/DynamicInput";
import {IssueCertificateRequest} from "./model/IssueCertificateRequest";
import ModalPreloader from "../../core/components/preloader/ModalPreloader";
import {IssueCertificateResponse} from "./model/IssueCertificateResponse";
import Notification from "../../core/components/notification/Notification";
import {UnexpectedResponseError} from "../../core/apiclient/ApiResponse";
import ValidationResultConverter from "../../core/validation/ValidationResultConverter";
import DomainNamesList from "./components/DomainNamesList";
import {Navigate} from "react-router-dom";

interface CertificateFormPageState {
    availableProviders: AvailableProviderResponse[]
    loading: boolean
    validationResult: ValidationResult
    formValues: IssueCertificateRequest
    certificateId?: string
}

export default class CertificateFormPage extends ShellAwareComponent<unknown, CertificateFormPageState> {
    private readonly service: CertificateService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance>

    constructor(props: any) {
        super(props);
        this.service = new CertificateService()
        this.saveModal = new ModalPreloader()
        this.formRef = React.createRef()
        this.state = {
            loading: true,
            validationResult: new ValidationResult(),
            formValues: {
                providerId: "",
                domainNames: [""],
                parameters: {}
            },
            availableProviders: [],
        }
    }

    private submit() {
        this.saveModal.show(
            "Hang on tight",
            "We're issuing your certificate. This can take several seconds, feel free to grab a cup of coffee."
        )

        const certificateRequest = { ...this.state.formValues }
        certificateRequest.domainNames = certificateRequest.domainNames.map(item => item ?? "")

        this.service
            .issue(certificateRequest)
            .then(response => this.handleResponse(response))
            .catch(error => this.handleError(error))
            .then(() => this.saveModal.close())
    }

    private handleResponse(response: IssueCertificateResponse) {
        const {success, errorReason, certificateId} = response
        if (success) {
            Notification.success(
                "Certificate issued",
                "The SSL certificate was issued and is now ready to be used"
            )
            this.setState({ certificateId })
        } else {
            Notification.error(
                "Issue failed",
                errorReason ?? "Unable to issue the certificate at this moment. Please try again later.",
            )
        }
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null)
                this.setState({ validationResult })
        }

        Notification.error(
            "That didn't work",
            "Please check the form to see if everything seems correct",
        )
    }

    private buildProviderSelectOptions() {
        const {availableProviders} = this.state
        return availableProviders.map(provider => ({
            value: provider.uniqueId,
            label: provider.name,
        }))
    }

    private updateFormState(formValues: IssueCertificateRequest) {
        this.setState({ formValues })
    }

    private renderDynamicFields() {
        const {formValues, availableProviders, validationResult} = this.state
        const provider = availableProviders.find(item => item.uniqueId === formValues.providerId)
        return provider
            ?.dynamicFields
            .sort((left, right) => left.priority > right.priority ? 1 : -1)
            .map(field => (
                <DynamicInput
                    validationResult={validationResult}
                    formValues={formValues}
                    field={field}
                />
            ))
    }

    private renderForm() {
        const {validationResult, formValues} = this.state

        return (
            <Form<IssueCertificateRequest>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                onValuesChange={(_, values) => this.updateFormState(values)}
                initialValues={formValues}
            >
                <Form.Item
                    name="providerId"
                    validateStatus={validationResult.getStatus("providerId")}
                    help={validationResult.getMessage("providerId")}
                    label="Certificate provider"
                    required
                >
                    <Select
                        placeholder="Certificate provider"
                        options={this.buildProviderSelectOptions()}
                    />
                </Form.Item>
                <DomainNamesList validationResult={validationResult} />
                {this.renderDynamicFields()}
            </Form>
        )
    }

    componentDidMount() {
        this.service
            .availableProviders()
            .then(providers => {
                this.setState({
                    availableProviders: providers,
                    loading: false,
                    formValues: {
                        providerId: providers[0].uniqueId,
                        domainNames: [""],
                        parameters: {},
                    }
                })
            })
    }

    shellConfig(): ShellConfig {
        return {
            title: "New SSL certificate",
            subtitle: "Issue or upload a SSL certificate for use with the nginx's virtual hosts",
            actions: [
                {
                    description: "Issue and save",
                    onClick: () => this.submit(),
                },
            ],
        }
    }

    render() {
        const {loading, availableProviders, certificateId} = this.state

        if (certificateId !== undefined)
            return <Navigate to={`/certificates/${certificateId}`} />

        return (
            <Preloader loading={loading}>
                <If condition={availableProviders.length > 0}>
                    {this.renderForm()}
                </If>
            </Preloader>
        )
    }
}

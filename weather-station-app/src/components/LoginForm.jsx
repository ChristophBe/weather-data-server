import * as React from "react";
import * as yup from 'yup';
import {Field, Form, Formik} from "formik/dist/index";
import PropTypes from "prop-types";
import {connect} from "react-redux";
import {InputField} from "./FormInputFields";

import AuthActions from "../redux/actions/AuthActions"
class  LoginForm extends React.Component {
    initialValue = { email: '', password: '', rememberMe:true };

    validationSchema =  yup.object().shape({
        email: yup.string().required("Bitte geben sie Ihre E-Mail Addresse ein.").email("Bitte geben sie eine Gültige E-Mail Addresse ein."),
        password: yup.string().required("Bitte geben sie Ihr Passwort ein.")
    });


    constructor(props) {
        super(props);
        this.onSubmitForm = this.onSubmitForm.bind(this);
    }


    onSubmitForm(values, { setSubmitting, resetForm }){

        console.log(values);
        this.props.doLogin(values);
        resetForm();
        setSubmitting(false);
    }

    render() {

        const {auth:{isAuthenticationFailed, isLoading}} = this.props;
        return (
            <div>
                {isAuthenticationFailed?<div className={"alert alert-danger"}>Die E-Mail Adresse oder das Passwort sind Fehlerhaft!</div> : null}
                <Formik
                    initialValues={this.initialValue}
                    validate={this.validateForm}
                    validationSchema={this.validationSchema}
                    onSubmit={this.onSubmitForm}>
                    {() => (
                        <Form>
                            <Field name="email" component={InputField} type={"email"} placeholder={"Ihre E-Mail Adresse"} label={"E-Mail Adresse"}/>
                            <Field name="password" component={InputField} type={"password"} placeholder={"Ihr Passwort"} label={"Passwort"}/>
                            <button className={"btn btn-outline-light"} type="submit" disabled={isLoading}>
                                {isLoading ? <span className="spinner-border spinner-border-sm mr-1" role="status" aria-hidden="true"/> : null}
                                Anmelden
                            </button>
                        </Form>
                    )}
                </Formik>
            </div>
        );
    }
}


LoginForm.propType= {
    setAuthTokenHandler: PropTypes.func
};


const mapStateToProps = (state) => ({
    auth:state.auth
});

const mapDispatchToProps = {
    doLogin: AuthActions.doLogin
};

export default connect(mapStateToProps,mapDispatchToProps)(LoginForm);
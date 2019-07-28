import * as React from "react";
import * as yup from 'yup';
import {Field, Form, Formik} from "formik/dist/index";
import {connect} from "react-redux";
import {InputField} from "./FormInputFields";
    import {NodesService} from "../services/NodesService";
import {UserService} from "../services/UserService";
class  CreateUserView extends React.Component {
    initialValue = { username: '' ,password:"", passwordConfirmation: ""};

    validationSchema =  yup.object().shape({
        username: yup.string().required("Bitte gib ein Nuternamen an.")
            .min(4, "Der Nutzername muss mindestens 4 Zeichen lang sein."),
        password: yup.string().required("Bitte gib ein Passwort an.")
            .min(4, "Das Passwort muss mindestens 4 Zeichen lang sein."),
        passwordConfirmation:  yup.string().oneOf([yup.ref('password'), null], 'Die Passwort wiederhohlung muss mit dem Passwort übereinstimmen.')
    });


    constructor(props) {
        super(props);

        this.state = {
            success: false,
            failed: false
        };
        this.onSubmitForm = this.onSubmitForm.bind(this);
    }


    onSubmitForm(values, { setSubmitting, resetForm }){

        console.log(values);
        const {token,auth} = this.props;

        const payload = {
            token,
            username: values.username,
            password: values.password
        }

        let reqPromise = UserService.createUser(payload);
        reqPromise.then(
            (data)=>{
                console.log("shareData",data)
                setSubmitting(false);
                this.setState({
                    success: true,
                    failed:false
                });
            },
            (err)=>{

                console.log("shareErr",err)
                setSubmitting(false);
                this.setState({
                    success: false,
                    failed:true
                });
            }
        );
        resetForm();
        setSubmitting(false);
    }

    render() {

        const {success,failed} = this.state;
        if(success){
            return <h2>Dein Benutzerkonto wurde Erfolgreich angelegt.</h2>
        }
        return (
            <div>
                <h1>Benutzerkonto anlegen</h1>
                {failed?<div className={"alert alert-danger"}>Es ist ein Fehler aufgetreten.</div> : null}
                <Formik
                    initialValues={this.initialValue}
                    validate={this.validateForm}
                    validationSchema={this.validationSchema}
                    onSubmit={this.onSubmitForm}>
                    {({isSubmitting}) => (
                        <Form>
                            <Field name="username" component={InputField} type={"text"} placeholder={"Benutzername"} label={"Benutzername"}/>
                            <Field name="password" component={InputField} type={"password"} placeholder={"Passwort"} label={"Passwort"}/>
                            <Field name="passwordConfirmation" component={InputField} type={"password"} placeholder={"Passwort Wiederhohlung"} label={"Passwort Wiederhohlung"}/>
                            <button className={"btn btn-outline-light"} type="submit" disabled={isSubmitting}>
                                {isSubmitting ? <span className="spinner-border spinner-border-sm mr-1" role="status" aria-hidden="true"/> : null}
                                anlegen
                            </button>
                        </Form>
                    )}
                </Formik>
            </div>
        );
    }
}


const mapStateToProps = (state) => ({
    auth:state.auth,
    nodes:state.nodes
});

const mapDispatchToProps = {};

export default connect(mapStateToProps,mapDispatchToProps)(CreateUserView);
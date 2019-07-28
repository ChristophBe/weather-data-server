import * as React from "react";
import * as yup from 'yup';
import {Field, Form, Formik} from "formik/dist/index";
import PropTypes from "prop-types";
import {connect} from "react-redux";
import {InputField} from "./FormInputFields";

import AuthActions from "../redux/actions/AuthActions"
import {NodesService} from "../services/NodesService";
class  ShareNodeView extends React.Component {
    initialValue = { email: '' };

    validationSchema =  yup.object().shape({
        email: yup.string().required("Bitte geben sie Ihre E-Mail Addresse ein.").email("Bitte geben sie eine Gültige E-Mail Addresse ein."),
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
        const {nodeId,auth} = this.props;

        let reqPromise = NodesService.shareNode(nodeId, values.email,auth.token);
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
            return <h2>Die Wetterstation wurde erfolgreich freigegeben.</h2>
        }

        const {nodes , nodeId} = this.props;
        const node = nodes.map[nodeId];
        console.log(node);
        return (
            <div>
                <h1>Wetterstation {node ? "\"" +node.name+  "\" ": ""}Freigeben</h1>
                {failed?<div className={"alert alert-danger"}>Es ist ein Fehler aufgetreten.</div> : null}
                <Formik
                    initialValues={this.initialValue}
                    validate={this.validateForm}
                    validationSchema={this.validationSchema}
                    onSubmit={this.onSubmitForm}>
                    {({isSubmitting}) => (
                        <Form>
                            <Field name="email" component={InputField} type={"email"} placeholder={"E-Mail Adresse"} label={"E-Mail Adresse"}/>
                            <button className={"btn btn-outline-light"} type="submit" disabled={isSubmitting}>
                                {isSubmitting ? <span className="spinner-border spinner-border-sm mr-1" role="status" aria-hidden="true"/> : null}
                                Teilen
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

export default connect(mapStateToProps,mapDispatchToProps)(ShareNodeView);
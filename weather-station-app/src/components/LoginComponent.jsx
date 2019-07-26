import * as React from "react";
import LoginForm from "./LoginForm";
import AuthActions from "../redux/actions/AuthActions";
import {connect} from "react-redux";

class LoginComponent extends React.Component{

    componentDidUpdate(prevProps, prevState, snapshot) {
        const {auth}= this.props;

        if(auth.token !== null){
            this.props.history.push('/')
        }
    }

    render() {


         return <div>
             <h2>Login</h2>

             <LoginForm/>

         </div>
     }
}


const mapStateToProps = (state) => ({
    auth:state.auth
});

const mapDispatchToProps = {
    doLogin: AuthActions.doLogin
};

export default connect(mapStateToProps,mapDispatchToProps)(LoginComponent);
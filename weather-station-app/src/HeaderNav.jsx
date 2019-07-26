import * as React from "react";

import {Link} from "react-router-dom"
import AuthActions from "./redux/actions/AuthActions";
import {connect} from "react-redux";

class HeaderNav  extends React.Component{


    constructor(props, context) {
        super(props, context);

        this.onLogoutButtonClicked = this.onLogoutButtonClicked.bind(this)
    }

    onLogoutButtonClicked(){
        this.props.doLogout();

    }
    render() {
        const {auth} = this.props
        return <nav className={"main-nav"}>
            <div className={"container"}>
                <nav>
                    <ul>
                        <li>

                        </li>

                    </ul>
                    <ul>
                        {
                            auth.currentUser != null ? <li>Hallo, <b>{auth.currentUser.username}</b></li> : null
                        }
                        <li>
                            {
                                auth.token == null ? <Link className={"btn btn-outline-light"} to={"/users/login"}>Login</Link> : <button onClick={this.onLogoutButtonClicked} className={"btn btn-outline-light"}>Logout</button>
                            }
                        </li>

                    </ul>
                </nav>
            </div>

        </nav>
    }
}


const mapStateToProps = (state) => ({
    auth:state.auth
});

const mapDispatchToProps = {
    doLogin: AuthActions.doLogin,
    doLogout: AuthActions.doLogout,
}

export default connect(mapStateToProps,mapDispatchToProps)(HeaderNav);
import React, {Component, Fragment} from 'react';
import chroma from 'chroma-js';
import HeaderNav from "./HeaderNav";

// eslint-disable-next-line no-unused-vars
import {BrowserRouter as Router, Route, Switch, Redirect, Link} from "react-router-dom";
import MeasuringComponent from "./components/MeasuringComponent";
import LoginComponent from "./components/LoginComponent";
import AuthActions from "./redux/actions/AuthActions";
import {connect} from "react-redux";
import NodesActions from "./redux/actions/NodesActions";



class App extends Component {
    constructor(props) {
        super(props);
    }

    componentDidMount(){
        this.props.tryCookieAuth();
        this.props.fetchNodes();
    }


    render() {

        console.log("props", this.props)
        const {nodes, measurements} = this.props;

        if(nodes.isLoading || nodes.selectedNode == null){
            return <div className="App" >
                <div className="main-content" >
                    <h1>Die Messstationen werden geladen</h1>
                </div>
            </div>
        }

        const style = {};

       if( measurements[nodes.selectedNode] != null
            && measurements[nodes.selectedNode].items.length >0){
            style["background"] = App.calcBackgroundColor(measurements[nodes.selectedNode].items[0].temperature)
        }




        return (
            <div className="App" style={style}>
                <Router>

                    <Fragment>
                        <HeaderNav/>
                        <section className="main-content" >

                            <Switch>
                                <Route path="/users/login" component={LoginComponent}/>

                                <Route path = "/node/:nodeId" render={({match})=><MeasuringComponent nodeId={Number(match.params.nodeId)}/>}/>
                                <Redirect from={"/"} to={"/node/" + nodes.selectedNode}/>
                            </Switch>
                        </section>
                        {Object.keys(nodes.map).length>1
                            ? <div className={"bottom-nav"}>
                                <div className={"container"}>
                                    <nav>
                                        {
                                            Object.keys(nodes.map).map((key,index)=><Link className={"btn btn-outline-light"} key={index} to={"/node/" + key}>{nodes.map[key].name}</Link>)
                                        }
                                    </nav>
                                </div>
                            </div>:null
                        }

                    </Fragment>
                </Router>

            </div>
         );
    }



    static calcBackgroundColor(currentTemperature) {
        //let scale = chroma.scale(['#007bff', '#fd7e14']);
        let scale = chroma.scale(['#0400ff','#00a1ff','#f9a100' ,  '#f94500']);
        //let scale = chroma.scale(['#66ccff', '#ffbe6e']);
        let tempOnColorScale = (currentTemperature + 10) / 50   ;
        tempOnColorScale = tempOnColorScale < 0 ? 0 : tempOnColorScale;
        tempOnColorScale = tempOnColorScale > 1 ? 1 : tempOnColorScale;


        return scale(tempOnColorScale).hex();
    }
}




const mapStateToProps = (state) => ({
    auth:state.auth,
    nodes:state.nodes,
    measurements:state.measurements,
});

const mapDispatchToProps = {
    doLogin: AuthActions.doLogin,
    tryCookieAuth: AuthActions.tryCookieAuth,
    fetchNodes: NodesActions.fechNodes
};

export default connect(mapStateToProps,mapDispatchToProps)(App);
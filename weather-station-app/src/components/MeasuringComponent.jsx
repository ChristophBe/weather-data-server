import * as React from "react";
import moment from "moment";
import MeasurementsActions from "../redux/actions/MeasurementsActions";
import {connect} from "react-redux";
import NodesActions from "../redux/actions/NodesActions";
import {Link} from "react-router-dom";
import {Fragment} from "react";

class MeasuringComponent extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            reloadInterval:false};
    }

   componentDidUpdate(prevProps, prevState, snapshot) {
        console.log("props test", this.props);
        if(prevProps.nodeId !== this.props.nodeId){

            this.props.fetch();
            this.props.select();
        }
   }

    componentDidMount() {
        this.props.fetch();
        this.props.select();
        let reloadInterval = setInterval(this.props.fetch,1000*60*2);
        this.setState({reloadInterval:reloadInterval})
    }

    componentWillUnmount(){
        if(this.state.reloadInterval){
            clearInterval(this.state.reloadInterval);
        }
    }



    render() {
        console.log("props a" , this.props);


        const {node,measurements, isOwnNode, nodeId} = this.props;

        if(measurements == null){
            return <h1>Messwerte werden geladen.</h1>
        }

        const{items, isLoading} = measurements;

        if(isLoading){
            return <h1>Messwerte werden geladen.</h1>
        }

        if(items.length <= 0){
            return <Fragment>
                <h1>Es gibt leider keine akutellen Messwerte.</h1>
                { node != null &&!node.is_public && isOwnNode ? <div><Link className={"btn btn-outline-light"} to={"/node/" + nodeId + "/share"}>Freigeben</Link></div>: null}
            </Fragment>
        }
        const lastMeasurement = items[0];

        const {temperature,humidity,pressure, timestamp} = lastMeasurement;

        const time = moment(timestamp);

        console.log("nodes Map"  , node);
        return (
            <div >
                {node !== undefined ? <h3>{node.name}</h3>: null}
                <h1>Temperatur: <b>{temperature}°C</b></h1>
                <div >
                    <p>Luftfeuchtigkeit: <b>{humidity}%</b></p>
                    <p>Luftdruck: <b>{pressure}hPa</b></p>
                    <p>Messzeitpunkt: <b>{time.format("DD.MM.YYYY HH:mm")} Uhr</b></p>
                </div>
                {!node.isPublic && isOwnNode ? <div><Link className={"btn btn-outline-light"} to={"/node/" + nodeId + "/share"}>Freigeben</Link></div>: null}
            </div>
        );
    }
}



const mapStateToProps = (state, ownProps) => ({
    node: state.nodes.map[ownProps.nodeId],
    isOwnNode: state.nodes.ownNodes.reduce(
        (accumulator, item) => { return accumulator || item.id === ownProps.nodeId}, false
    ),
    measurements: state.measurements[ownProps.nodeId],
});


const mapDispatchToProps = (dispatch,ownProps) => {
    return {
        fetch: () => {
            dispatch(MeasurementsActions.fechMeasurements(ownProps.nodeId));
        },
        select: () => {
            dispatch(NodesActions.selectNode(ownProps.nodeId));
        }
    };
};
export default connect(mapStateToProps,mapDispatchToProps)(MeasuringComponent);




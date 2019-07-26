import * as React from "react";
import {LineSeries, XAxis, XYPlot, YAxis} from "react-vis/es";

import {timeFormatDefaultLocale} from 'd3-time-format';

timeFormatDefaultLocale({
    "dateTime": "%A, der %e. %B %Y, %X",
    "date": "%d.%m.%Y",
    "time": "%H:%M:%S",
    "periods": [],
    "days": ["Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"],
    "shortDays": ["So", "Mo", "Di", "Mi", "Do", "Fr", "Sa"],
    "months": ["Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"],
    "shortMonths": ["Jan", "Feb", "Mrz", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Okt", "Nov", "Dez"]
});

class Graphs extends React.Component{
    render() {

        let temperatures =this.props.measurements.map((value) =>  ({x: new Date(value.timestamp),y: value.temperature}));
        let tempMinMax = this.generateMinMax(temperatures, 2);
        let humiditys =this.props.measurements.map((value) =>  ({x: new Date(value.timestamp),y: value.humidity}));
        let humiditysMinMax = this.generateMinMax(humiditys, 10);
        let pressures =this.props.measurements.map((value) =>  ({x: new Date(value.timestamp),y: value.pressure}));
        let pressuresMinMax = this.generateMinMax(pressures, 50);


        //console.log(temperatures);
        return (<div className={"container"}>
                        <div className={"row"}>

                            <XYPlot height={200} width={400} className={"col-4"} yDomain={tempMinMax} xType="time" yType="linear">
                                <XAxis title="Zeit" />
                                <YAxis title="Temeratur °C" />
                                <LineSeries data={temperatures} style={{strokeWidth: 2, fill:0  }}/>
                            </XYPlot>
                            <XYPlot height={200} width={400} className={"col-4"} yDomain={humiditysMinMax} xType="time" yType="linear">
                                <XAxis title="Zeit" />
                                <YAxis title="Luftfeuchtigkeit %" />
                                <LineSeries data={humiditys} style={{strokeWidth: 2, fill:0  }}/>
                            </XYPlot>
                            <XYPlot height={200} width={400} className={"col-4"} yDomain={pressuresMinMax} xType="time" yType="linear">
                                <XAxis title="Zeit" />
                                <YAxis title="Luftdruck hPa" />
                                <LineSeries data={pressures} style={{strokeWidth: 2, fill:0  }}/>
                            </XYPlot>
                        </div>
                    </div>
        );
    }

    generateMinMax(temperatures, diff) {
        let min = temperatures[0].y;
        let max = temperatures[0].y;
        temperatures.forEach(item => {
            if (item.y < min) {
                min = item.y
            }
            if (item.y > max) {
                max = item.y
            }
        });
        let tempMinMax = [min - diff, max + diff];
        return tempMinMax;
    }

}
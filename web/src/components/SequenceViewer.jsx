import React from 'react';

export default class SequenceViewer extends React.Component {
    componentDidMount() {
        var FeatureViewer = require("feature-viewer");
        var ft2 = new FeatureViewer("FDSJKLFJDSFKLJDFHADJKLFHDSJKLFHDAFJKLDHFJKLDASFHDJKLFHDSAJKLFHDAKLFJDHSAFKLDLSNCDJKLFENFIUPERWDJKPCNVDFPIEHFDCFJDKOWFPDJWFKLXSJFDW9FIPUAENDCXAMSFNDUAFIDJFDLKSAFJDSAKFLJDSADJFDW9FIPUAENDCXAMSFNDAAAAAAAAAAAFJDSAKFL", "#seq", {
            showAxis: true,
            showSequence: true,
            brushActive: true,
            toolbar: true,
            bubbleHelp: true,
            zoomMax: 10
        });

        ft2.addFeature({
            data: [{ x: 20, y: 40 }, { x: 46, y: 100 }, { x: 123, y: 167 }],
            name: "1B93",
            className: "test1",
            color: "#005572",
            type: "rect",
            filter: "type1"
        });
        ft2.addFeature({
            data: [{ x: 50, y: 90, description: "PF666" }],
            name: "Pfam",
            className: "test6",
            color: "#81BEAA",
            type: "rect",
            filter: "type2"
        });
        ft2.addFeature({
            data: [{ x: 52, y: 52 }, { x: 92, y: 92 }],
            name: "Variations",
            className: "test2",
            color: "#006588",
            type: "rect",
            filter: "type2"
        });
    }
    render() {
        return (
            <div id="seq" />
        )
    }
}
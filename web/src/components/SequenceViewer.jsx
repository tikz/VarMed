import React from 'react';

export default class SequenceViewer extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            pos: 0
        };

        this.handleMouseMove = this.handleMouseMove.bind(this);
        this.handleMouseLeave = this.handleMouseLeave.bind(this);
    }

    // This is hackish and adds overhead, but FeatureViewer doesn't provide mousemove or
    // position change events and very little in the API is exported. In the future, maybe
    // ask the authors on GitHub for this feature, or fork the repo and add it.
    handleMouseMove(e) {
        var posText = this.state.zoomPositionElement.innerText;
        if (posText != this.state.pos) {
            var pos = parseInt(posText.slice(0, posText.length - 1));
            this.props.highlight(pos, pos);
            this.setState(state => ({
                pos: posText
            }));
        }
    };

    handleMouseLeave(e) {
        this.props.highlight(0, 0);
    };

    componentDidMount() {
        var FeatureViewer = require("feature-viewer");
        var fv = new FeatureViewer("MTEYKLVVVGAGGVGKSALTIQLIQNHFVDEYDPTIEDSYRKQVVIDGETCLLDILDTAGQEEYSAMRDQYMRTGEGFLCVFAINNSKSFADINLYREQIKRVKDSDDVPMVLVGNKCDLPTRTVDTKQAHELAKSYGIPFIETSAKTRQGVEDAFYTLVREIRQYRMKKLNSSDDGTQGCMGLPCVVM", "#fv", {
            showAxis: true,
            showSequence: true,
            brushActive: true,
            toolbar: true,
            bubbleHelp: false,
            zoomMax: 20
        });

        fv.addFeature({
            data: [{ x: 1, y: 168, description: "Chain A" }],
            name: "3CON",
            className: "test1",
            color: "#2196F3",
            type: "rect",
            filter: "type1"
        });
        fv.addFeature({
            data: [{ x: 5, y: 165, description: "PF00071 Ras" }],
            name: "Pfam",
            className: "test6",
            color: "#2196F3",
            type: "rect",
            filter: "type2"
        });
        fv.addFeature({
            data: [{ x: 52, y: 52 }, { x: 92, y: 92 }],
            name: "Variant",
            className: "variant",
            color: "#21CBF3",
            type: "rect",
            filter: "type3"
        });

        var selectFunc = this.props.select
        fv.onFeatureSelected(function (d) {
            console.log(d);
            selectFunc(d.detail.start, d.detail.end);
        });

        fv.onZoom(function (d) {
            console.log(d.detail);
        });

        this.setState(state => ({
            zoomPositionElement: document.getElementById("zoomPosition")
        }));

    };
    render() {
        return (
            <div id="fv" onMouseMove={this.handleMouseMove} onMouseLeave={this.handleMouseLeave} />
        )
    };
}
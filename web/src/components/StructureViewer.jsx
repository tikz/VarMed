import React from 'react';
import LiteMol from 'litemol';
import 'litemol/dist/css/LiteMol-plugin.css';

export default class StructureViewer extends React.Component {
    componentDidMount() {
        var plugin = LiteMol.Plugin.create({
            target: '#litemol',
            layoutState: {
                hideControls: true,
                isExpanded: false
            },
        });
        plugin.loadMolecule({
            id: '1tqn',
            url: 'https://www.ebi.ac.uk/pdbe/static/entry/1tqn_updated.cif',
            format: 'cif'
        });
    }
    render() {
        return (
            <div id="litemol" style={{ position: 'relative', height: "600px" }}></div>
        )
    }
}

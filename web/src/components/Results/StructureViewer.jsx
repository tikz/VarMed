import React from 'react';
import LiteMol from 'litemol';
import 'litemol/dist/css/LiteMol-plugin.css';


let Transformer = LiteMol.Bootstrap.Entity.Transformer;
let Transform = LiteMol.Bootstrap.Tree.Transform;

export default class StructureViewer extends React.Component {
    highlight(start, end) {
        this.clearHighlight();
        var plugin = this.state.plugin
        var model = plugin.context.select('model')[0];
        let query = LiteMol.Core.Structure.Query.sequence('1', 'A', { seqNumber: start }, { seqNumber: end });
        LiteMol.Bootstrap.Command.Molecule.Highlight.dispatch(plugin.context, { model, query, isOn: true });
    };

    clearHighlight() {
        var plugin = this.state.plugin
        var model = plugin.context.select('model')[0];
        let query = LiteMol.Core.Structure.Query.sequence('1', 'A', { seqNumber: 1 }, { seqNumber: 185 });
        LiteMol.Bootstrap.Command.Molecule.Highlight.dispatch(plugin.context, { model, query, isOn: false });
    };

    focus(start, end) {
        var plugin = this.state.plugin
        var model = plugin.context.select('model')[0];
        let query = LiteMol.Core.Structure.Query.sequence('1', 'A', { seqNumber: start }, { seqNumber: end });
        LiteMol.Bootstrap.Command.Molecule.FocusQuery.dispatch(plugin.context, { model, query, isOn: true });
    };

    select(start, end) {
        var plugin = this.state.plugin
        var model = plugin.context.select('model')[0];
        let query = LiteMol.Core.Structure.Query.sequence('1', 'A', { seqNumber: start }, { seqNumber: end });
        LiteMol.Bootstrap.Command.Molecule.CreateSelectInteraction.dispatch(plugin.context, { entity: model, query: query });
    };

    applyTheme() {
        var plugin = this.state.plugin;
        let colors = LiteMol.Core.Utils.FastMap.create();
        // colors.set('Uniform', LiteMol.Visualization.Color.fromHex(0x095c64));
        // colors.set('Uniform', LiteMol.Visualization.Color.fromHex(0x006e70));
        colors.set('Uniform', LiteMol.Visualization.Color.fromHex(0x006870));
        colors.set('Selection', LiteMol.Visualization.Color.fromHex(0xf15a29));
        colors.set('Highlight', LiteMol.Visualization.Color.fromHex(0xff8a2b));
        let theme = LiteMol.Bootstrap.Visualization.Molecule.uniformThemeProvider(void 0, { colors });

        const visuals = plugin.selectEntities(LiteMol.Bootstrap.Tree.Selection.byRef('polymer-visual').subtree().ofType(LiteMol.Bootstrap.Entity.Molecule.Visual));
        for (const v of visuals) {
            plugin.command(LiteMol.Bootstrap.Command.Visual.UpdateBasicTheme, { visual: v, theme });
        }
    };

    componentDidMount() {
        var plugin = LiteMol.Plugin.create({
            target: '#litemol',
            viewportBackground: '#1c1e20',
            layoutState: {
                hideControls: true,
                isExpanded: false
            },
        });

        let id = '3con';
        let action = Transform.build()
            .add(plugin.context.tree.root, Transformer.Data.Download, { url: `https://www.ebi.ac.uk/pdbe/static/entry/${id}_updated.cif`, type: 'String', id })
            .then(Transformer.Data.ParseCif, { id }, { isBinding: true })
            .then(Transformer.Molecule.CreateFromMmCif, { blockIndex: 0 }, { isBinding: true })
            .then(Transformer.Molecule.CreateModel, { modelIndex: 0 }, { isBinding: false, ref: 'model' })
            .then(Transformer.Molecule.CreateMacromoleculeVisual, { polymer: true, polymerRef: 'polymer-visual', het: true, water: true });

        plugin.applyTransform(action).then(() => {
            this.setState(state => ({
                plugin: plugin
            }));
            this.applyTheme();
        });
    };

    render() {
        return (
            <div id="litemol" />
        )
    };
}

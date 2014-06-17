Ext.define('IW.view.wiki.Window', {
    extend: 'Ext.window.Window',
    xtype: 'basic-window',
    alias: 'widget.wikiwindow',

    title: 'Wikis',
    height: 400,
    width: 400,
    layout: 'fit',
    autoShow: true,
    collapsible: true,
    closable: false,

    tools: [{
        type: 'plus',
        tooltip: 'Create new Wiki',
        width: 16,
        handler: function(event, target, owner, tool) {
            this.up('window').fireEvent('itemnewbuttonclick');
        }
    }],

    collapsible: true,
    initComponent: function() {

        this.items = [
            {
                xtype: 'wikilist'
            }
        ];

        this.callParent(arguments);
    }
});


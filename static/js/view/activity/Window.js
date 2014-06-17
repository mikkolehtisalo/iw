Ext.define('IW.view.activity.Window', {
    extend: 'Ext.window.Window',
    xtype: 'basic-window',
    alias: 'widget.activitywindow',

    title: 'Activities',
    height: 400,
    width: 400,
    layout: 'fit',
    autoShow: true,
    collapsible: true,
    closable: false,

    collapsible: true,
    initComponent: function() {

        this.items = [
            {
                xtype: 'activitylist'
            }
        ];

        this.callParent(arguments);
    }
});


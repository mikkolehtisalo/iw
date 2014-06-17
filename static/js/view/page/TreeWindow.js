Ext.define('IW.view.page.TreeWindow', {
    extend: 'Ext.window.Window',
    xtype: 'treewindow',
    alias: 'widget.treewindow',

    title: 'Pages',
    height: 400,
    width: 400,
    autoShow: true,
    overflowY: 'auto',
    collapsible: true,
    tools: [{
        type: 'Refresh',
        tooltip: 'Refresh',
        width: 16,
        handler: function(event, target, owner, tool) {
            this.up('window').fireEvent('refreshtree',event,target,owner,tool);
        }
    },{
        type: 'Plus',
        tooltip: 'Add',
        width: 16,
        handler: function(event, target, owner, tool) {
            this.up('window').fireEvent('addpagebuttonclick',event,target,owner,tool);
        }
    },{
        type: 'Minus',
        tooltip: 'Delete',
        hidden: true,
        width: 16,
        handler: function(event, target, owner, tool) {
            this.up('window').fireEvent('deletepagebuttonclick',event,target,owner,tool);
        }
    },{
        type: 'Key',
        tooltip: 'Access',
        hidden: true,
        width: 12,
        handler: function(event, target, owner, tool) {
            this.up('window').fireEvent('pageaccessrightsbuttonclick',event,target,owner,tool);
        }
    }
    ],

    initComponent: function() {

        this.items = [
            {
                xtype: 'pagetree',
                record: this.record
            }
        ];

        windowCount=1;
        Ext.WindowManager.each(
            function(item) {
                if (item.xtype == 'treewindow') {
                    windowCount ++;
                }
            });

        this.x = (450 + (windowCount * 25)) % 1300;
        this.y = (20 + (windowCount * 25)) % 800;

        // Set the title depending on the selected wiki
        this.title = this.record.data.Title;
        this.callParent(arguments);
    }
});

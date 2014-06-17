
Ext.define('IW.view.page.PageWindow', {
    extend: 'Ext.window.Window',
    alias: 'widget.pagewindow',
    xtype: 'pagewindow',
    requires: [
        'IW.Utilities'
    ],

    height: 600,
    width: 600,
    title: 'Window',
    collapsible: true,
    layout: 'fit',
    tools: [{
        type: 'Edit',
        tooltip: 'Edit',
        width: 16,
        handler: function(event, target, owner, tool) {
            this.up('window').fireEvent('itemeditbuttonclick',event,target,owner,tool);
        }
    },{
        type: 'Refresh',
        tooltip: 'Refresh',
        width: 16,
        handler: function(event, target, owner, tool) {
            this.up('window').fireEvent('itemrefreshbuttonclick',event,target,owner,tool);
        }
    },{
        type: 'Save',
        tooltip: 'Save',
        hidden: true,
        width: 16,
        margin: '0 1px 0 4px',
        handler: function(event, target, owner, tool) {
            this.up('window').fireEvent('itemsavebuttonclick',event,target,owner,tool);
        }
    },{
        type: 'Cancel',
        tooltip: 'Cancel',
        hidden: true,
        width: 16,
                margin: '0 1px 0 1px',
        handler: function(event, target, owner, tool) {
            this.up('window').fireEvent('itemcancelbuttonclick',event,target,owner,tool);
        }
    }],

    initComponent: function() {
        this.items = [{
        xtype: 'panel',
        bodyCls: 'pagecontent',
        border: 0,
        autoScroll: true
    }];

        windowCount=1;
        Ext.WindowManager.each(
            function(item) {
                if (item.xtype == 'pagewindow') {
                    windowCount ++;
                }
                
            });

        this.x = (900 + (windowCount * 25)) % 1300;
        this.y = (20 + (windowCount * 25)) % 800;

        this.callParent(arguments);

        if (!IW.Utilities.canwrite(this.record)) {
            this.tools[0].hidden = true; // Hide the edit button if we can't edit the content anyways
        }

    }

});


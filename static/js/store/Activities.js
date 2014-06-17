Ext.define('IW.store.Activities', {
    extend: 'Ext.data.Store',
    model: 'IW.model.Activity',
    autoLoad: false,
    proxy: {
        type: 'rest',
        url: '/api/activities/',
        reader: {
            type: 'json',
            successProperty: 'success'
        }
    },
    sorters: [{
        property: 'Timestamp',
        direction: 'DESC'
    }],
    constructor: function() {
        this.callParent(arguments);

        // Start automatic refresh
        var runner = new Ext.util.TaskRunner();
        var me = this;
        this.refreshtask = runner.start({
            run: function() {
                if (me.getCount()>0) {
                    if (sessionStorage.iw_csrf_token) {
                        me.reload();
                    }
                }
                
            },
            interval: 300000 // Once per 5 minutes should do
        });
    },
    destroy: function() {
        this.refreshtask.destroy();
        this.callParent(arguments);
    }
});



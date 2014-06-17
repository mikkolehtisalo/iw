Ext.define('IW.store.Wikis', {
    extend: 'Ext.data.Store',
    model: 'IW.model.Wiki',
    autoLoad: false,
    sorters: [{
        property: 'favorite',
        direction: 'DESC'
    },{
        property: 'title',
        direction: 'ASC'
    }],
    proxy: {
        type: 'rest',
        url: '/api/wikis',
        reader: {
            type: 'json',
            root: 'wikis',
            successProperty: 'success'
        }
    },
    constructor: function() {
        this.callParent(arguments);

        // Start automatic list refresh
        var runner = new Ext.util.TaskRunner();
        var me = this;
        this.refreshtask = runner.start({
            run: function() {
                if (me.getCount()>0) {
                    me.reload();
                }
                
            },
            interval: 300000 // Once per 5 minutes should do
        });
    },
    destroy: function() {
        this.refreshtask.destroy();
        this.callParent(arguments);
    },
    listeners: {
            'load' :  function(store,records,options) {
                // The Wikis store was loaded, so let's load FavoriteWikis
                var fw = Ext.data.StoreManager.lookup('FavoriteWikis');
                fw.load();
            }
    }
});


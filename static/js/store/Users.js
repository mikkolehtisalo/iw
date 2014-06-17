Ext.define('IW.store.Users', {
    extend: 'Ext.data.Store',
    model: 'IW.model.User',
    autoLoad: false,
    proxy: {
        type: 'rest',
        url: '/api/users/',
        reader: {
            type: 'json',
            successProperty: 'success'
        }
    },
    constructor: function() {
        // Set up the csrf token if available
        this.proxy.extraParams = {
            csrf_token: sessionStorage.iw_csrf_token
        };
        this.callParent(arguments);
    }
});


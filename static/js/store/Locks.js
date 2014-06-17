Ext.define('IW.store.Locks', {
    extend: 'Ext.data.Store',
    model: 'IW.model.Lock',
    autoLoad: false,
    proxy: {
        type: 'rest',
        url: '/api/locks/',
        reader: {
            type: 'json',
            successProperty: 'success'
        }
    }
});


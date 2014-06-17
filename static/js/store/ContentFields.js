Ext.define('IW.store.ContentFields', {
    extend: 'Ext.data.Store',
    model: 'IW.model.ContentField',
    autoLoad: false,
    proxy: {
        type: 'rest',
        url: '/api/contentfields/',
        reader: {
            type: 'json',
            successProperty: 'success'
        }
    }
});


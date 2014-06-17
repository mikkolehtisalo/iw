Ext.define('IW.store.Pages', {
    extend: 'Ext.data.TreeStore',
    model: 'IW.model.Page',
    autoLoad: false,
    clearOnLoad: true,

    proxy: {
        type: 'rest',
        url: '/api/pages/',
        reader: {
            type: 'json',
        },
        writer: {
            type: 'json',
        }
    }
});

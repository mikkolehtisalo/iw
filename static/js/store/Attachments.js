Ext.define('IW.store.Attachments', {
    extend: 'Ext.data.Store',
    model: 'IW.model.Attachment',
    autoLoad: false,
    proxy: {
        type: 'rest',
        url: '/api/attachments/',
        reader: {
            type: 'json',
            successProperty: 'success'
        }
    }
});


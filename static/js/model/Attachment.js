Ext.define('IW.model.Attachment', {
    extend: 'Ext.data.Model',
    fields: [
        'Attachment_id', 
        'Wiki_id', 
        'Attachment', 
        'Mime',
        'Filename',
        'Modified',
        'Status'
    ],
    //idProperty: 'Attachment_id',
    proxy: {
        type: 'rest',
        url: '/api/attachments/',
        reader: {
            type: 'json',
            successProperty: 'success'
        }
    }
});


Ext.define('IW.model.Lock', {
    extend: 'Ext.data.Model',
    fields: [
        {
            name: 'id',
            type: 'string',
            convert: function(value, record) {
                return record.get('Wiki_id')+'/'+record.get('Target_id');
            }
        },
        'Target_id', 
        'Wiki_id', 
        'Username',
        'Realname',
        'Modified',
    ],
    idProperty: 'id',
    proxy: {
        type: 'rest',
        url: '/api/locks/',
        reader: {
            type: 'json',
            successProperty: 'success'
        }
    }
});



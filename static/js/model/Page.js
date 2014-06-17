Ext.define('IW.model.Page', {
    extend: 'Ext.data.Model',
    fields: [
        {
            name: 'id',
            type: 'string',
            convert: function(value, record) {
                return record.get('Wiki_id')+'/'+record.get('Page_id');
            }
        },
        'Page_id', 
        'Wiki_id', 
        'Path', 
        'Title', 
        'Create_user', 
        'Readacl',
        'Writeacl',
        'Adminacl',
        'Stopinheritation',
        'Modified',
        { name: 'MatchedPermissions', defaultValue: []},
        { 
            name: 'index', 
            type: 'int',  
            defaultValue: null, 
            persist: true 
        },
    ],
    idProperty: 'id'
});



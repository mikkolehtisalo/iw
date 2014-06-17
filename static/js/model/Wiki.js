Ext.define('IW.model.Wiki', {
    extend: 'Ext.data.Model',
    fields: [
    'Wiki_id', 
    'Title', 
    'Description', 
    'Create_user',
    'Readacl',
    'Writeacl',
    'Adminacl',
    { name: 'MatchedPermissions', defaultValue: []},
    { name: 'Favorite', persist : false, defaultValue: false},
    ],
    idgen: 'uuid',
    idProperty: 'Wiki_id'
});

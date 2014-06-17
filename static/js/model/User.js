Ext.define('IW.model.User', {
    extend: 'Ext.data.Model',
    fields: ['username', 'name', 'groups', 'image'],
    idProperty: 'username'
});

Ext.define('IW.store.UserGroupSearch', {
    extend: 'Ext.data.Store',
    model: 'IW.model.UserGroupSearchItem',
    autoLoad: false,
    proxy: {
        type: 'rest',
        url: '/api/usergroupsearch/',
        reader: {
            type: 'json',
            successProperty: 'success'
        }
    }
});


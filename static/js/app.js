
Ext.Loader.setPath('Ext', '/static/extjs/src/');

Ext.require('Ext.container.Viewport');
Ext.require('Ext.Date');

Ext.define('IW.Utilities', {
    statics: {
        canread: function (record) {
            var permissions = record.get('MatchedPermissions'); 
            return (permissions.indexOf('admin')!=-1)||(permissions.indexOf('write')!=-1)||(permissions.indexOf('read')!=-1);
        },
        canwrite: function (record) {
            var permissions = record.get('MatchedPermissions'); 
            return (permissions.indexOf('admin')!=-1)||(permissions.indexOf('write')!=-1);
        },
        canadmin: function (record) {
            var permissions = record.get('MatchedPermissions'); 
            return (permissions.indexOf('admin')!=-1);
        }
    }
});

Ext.application({
    requires: ['Ext.container.Viewport'],
    name: 'IW',
    appFolder: '/static/js',
    controllers: [
        'Wikis',
        'Pages',
        'Activities'
    ],
    launch: function() {

        // Every Ajax call should be include the csrf token as header
        Ext.Ajax.defaultHeaders = {
            'X-CSRF-Token' : sessionStorage.iw_csrf_token
        };

        // Load what is needed at first
        Ext.data.StoreManager.each(function (item, index, len) {
            if ((item.storeId == 'Wikis' || item.storeId == 'Activities'))  {
                item.load();
            }
        });

        // Create the UI
        Ext.create('Ext.container.Viewport', {
            layout: 'absolute',
            items: [
                {
                    xtype: 'wikiwindow',
                    x: 20,
                    y: 20
                },
                {
                    xtype: 'activitywindow',
                    x: 20,
                    y: 430
                }
            ]
        });

    }
});


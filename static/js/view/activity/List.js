
Ext.define('IW.view.activity.List' ,{
    extend: 'Ext.grid.Panel',
    alias: 'widget.activitylist',

    invalidateScrollerOnRefresh: false,
    viewConfig: {
        loadMask: false
    },
    hideHeaders: true,
    border: 0,

    store: 'Activities',
    getSince: function(dateString) {
        var dt = new Date(dateString);
        var seconds = Math.floor(Ext.Date.getElapsed(dt) / 1000);

        var interval = Math.floor(seconds / 31536000);

        if (interval > 1) {
            return interval + " years";
        }
        interval = Math.floor(seconds / 2592000);
        if (interval > 1) {
            return interval + " months";
        }
        interval = Math.floor(seconds / 86400);
        if (interval > 1) {
            return interval + " days";
        }
        interval = Math.floor(seconds / 3600);
        if (interval > 1) {
            return interval + " hours";
        }
        interval = Math.floor(seconds / 60);
        if (interval > 1) {
            return interval + " minutes";
        }
        return Math.floor(seconds) + " seconds";
    },
    initComponent: function() {
        this.columns = [
            {header: 'Image',
            width: 48,
            renderer: function (value, metaData, record, row, col, store, gridView) {
                return '<img src="/user/'+record.data.User_id+'.jpeg" height="42" width="42" />';
            }
            },
            {header: 'Name', 
             dataIndex: 'Activity_id',
             flex: 1,
             renderer: function (value, metaData, record, row, col, store, gridView) {
                var template = '{0} {1} {2} <em>{3}</em> {4} ago';
                var action = '';
                if (record.data.Activity_type == 'ACTIVE') {
                    action = 'modified';
                } else if (record.data.Activity_type == 'DELETED') {
                    action = 'deleted';
                }

                var since = this.getSince(record.data.Timestamp);

                return Ext.String.format(template, 
                    record.data.User_name, 
                    action, 
                    Ext.util.Format.lowercase(record.data.Target_type),
                    record.data.Target_title,
                    since
                    );
             }
            },
        ];

        this.callParent(arguments);
    }
});

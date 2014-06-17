Ext.define('IW.model.Activity', {
    extend: 'Ext.data.Model',
    fields: [
    'Activity_id', 
    'Timestamp', 
    'User_id',
    'User_name', 
    'Activity_type', 
    'Target_title',
    'Target_type',
    'Target_id'
    ],
    idProperty: 'Activity_id'
});


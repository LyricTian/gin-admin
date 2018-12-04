-- 创建数据库
CREATE DATABASE `ginadmin` DEFAULT CHARACTER SET = `utf8mb4`;

-- 创建`g_menu`表
CREATE TABLE `g_menu` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `record_id` varchar(36) DEFAULT NULL,
  `code` varchar(50) DEFAULT NULL,
  `name` varchar(50) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `sequence` int(11) DEFAULT NULL,
  `icon` varchar(200) DEFAULT NULL,
  `path` varchar(200) DEFAULT NULL,
  `method` varchar(50) DEFAULT NULL,
  `level_code` varchar(20) DEFAULT NULL,
  `parent_id` varchar(36) DEFAULT NULL,
  `is_hide` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `creator` varchar(36) DEFAULT NULL,
  `created` bigint(20) DEFAULT NULL,
  `updated` bigint(20) DEFAULT NULL,
  `deleted` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_record_id` (`record_id`),
  KEY `idx_code` (`code`),
  KEY `idx_name` (`name`),
  KEY `idx_type` (`type`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted` (`deleted`),
  KEY `idx_is_hide` (`is_hide`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

LOCK TABLES `g_menu` WRITE;

INSERT INTO `g_menu` (`id`, `record_id`, `code`, `name`, `type`, `sequence`, `icon`, `path`, `method`, `level_code`, `parent_id`, `is_hide`, `status`, `creator`, `created`, `updated`, `deleted`)
VALUES
	(15,'047aecdc-76c8-4bfd-8dbc-02a37295d40b','admin','权限管理',10,90,'','','','01','',2,1,'',1543546798,0,0),
	(16,'d1ef3f75-ebc1-4b0d-be69-25e406b843af','system','系统管理',20,90,'setting','','','0101','047aecdc-76c8-4bfd-8dbc-02a37295d40b',2,1,'',1543546817,1543932539,0),
	(17,'751ffa55-fcbb-43bc-8b63-c3287f1f42d6','menu','菜单管理',30,10,'solution','/system/menu','','010101','d1ef3f75-ebc1-4b0d-be69-25e406b843af',2,1,'',1543546836,1543932569,0),
	(18,'7f6c7556-5242-444f-9714-59a1b5d1abcf','role','角色管理',30,20,'audit','/system/role','','010102','d1ef3f75-ebc1-4b0d-be69-25e406b843af',2,1,'',1543546953,1543932702,0),
	(19,'4b3448fd-c23f-49df-a51b-94aa2a68aec6','user','用户管理',30,30,'user','/system/user','','010103','d1ef3f75-ebc1-4b0d-be69-25e406b843af',2,1,'',1543546963,1543932708,0),
	(25,'14f966de-a307-4731-bc9a-8889b2b5a1dd','query','查询菜单数据',40,1,'','/api/v1/menus','GET','01010101','751ffa55-fcbb-43bc-8b63-c3287f1f42d6',1,1,'root',1543927860,0,0),
	(26,'0851bc50-5225-423a-a189-54cee416737a','one','查询指定菜单数据',40,2,'','/api/v1/menus/:id','GET','01010102','751ffa55-fcbb-43bc-8b63-c3287f1f42d6',1,1,'root',1543927900,0,0),
	(27,'36a8350c-5cd6-45ed-9734-85f3c990a905','create','创建菜单数据',40,3,'','/api/v1/menus','POST','01010103','751ffa55-fcbb-43bc-8b63-c3287f1f42d6',1,1,'root',1543927924,0,0),
	(28,'6c925f7c-f949-4f87-91ba-8e28c3cf0ab3','update','更新菜单数据',40,4,'','/api/v1/menus/:id','PUT','01010104','751ffa55-fcbb-43bc-8b63-c3287f1f42d6',1,1,'root',1543928525,0,0),
	(29,'e3ba022a-060e-4b0d-af20-f416c333bdc3','delete','删除菜单数据',40,5,'','/api/v1/menus/:id','DELETE','01010105','751ffa55-fcbb-43bc-8b63-c3287f1f42d6',1,1,'root',1543928605,0,0),
	(30,'000b3c20-95e0-4d07-ab54-9e1a092d8b86','deleteMany','删除多条菜单数据',40,6,'','/api/v1/menus','DELETE','01010106','751ffa55-fcbb-43bc-8b63-c3287f1f42d6',1,1,'root',1543928642,0,0),
	(31,'82894ddb-0359-4b61-9796-e9af9631b053','enable','启用菜单数据',40,7,'','/api/v1/menus/:id/enable','PATCH','01010107','751ffa55-fcbb-43bc-8b63-c3287f1f42d6',1,1,'root',1543928710,0,0),
	(32,'4be7cfe2-e7b3-4e52-b76f-0d19dec7aad5','disable','禁用菜单数据',40,8,'','/api/v1/menus/:id/disable','PATCH','01010108','751ffa55-fcbb-43bc-8b63-c3287f1f42d6',1,1,'root',1543928752,0,0),
	(33,'af4edcc6-28fd-4aad-b67a-1653c2a7d0e2','query','查询角色数据',40,1,'','/api/v1/roles','GET','01010201','7f6c7556-5242-444f-9714-59a1b5d1abcf',1,1,'root',1543932205,0,0),
	(34,'3c039cae-2769-476f-890f-183f4effc987','one','查询指定角色数据',40,2,'','/api/v1/roles/:id','GET','01010202','7f6c7556-5242-444f-9714-59a1b5d1abcf',1,1,'root',1543932224,0,0),
	(35,'88555d41-f564-45e8-bdad-ec324915d124','create','创建角色数据',40,3,'','/api/v1/roles','POST','01010203','7f6c7556-5242-444f-9714-59a1b5d1abcf',1,1,'root',1543932247,0,0),
	(36,'a8208bb8-4a2b-45e5-a53d-5d4fe9b7f8a8','update','更新角色数据',40,4,'','/api/v1/roles/:id','PUT','01010204','7f6c7556-5242-444f-9714-59a1b5d1abcf',1,1,'root',1543932268,0,0),
	(37,'d4ca0149-4475-4d63-be43-d9ae7e7794a5','delete','删除角色数据',40,5,'','/api/v1/roles/:id','DELETE','01010205','7f6c7556-5242-444f-9714-59a1b5d1abcf',1,1,'root',1543932292,0,0),
	(38,'37c7caaf-bd68-4a22-995f-be81cd8964b0','deleteMany','删除多条数据数据',40,6,'','/api/v1/roles','DELETE','01010206','7f6c7556-5242-444f-9714-59a1b5d1abcf',1,1,'root',1543932316,0,0),
	(39,'5755257d-cd9f-4db2-a874-0d88622e8e8e','enable','启用角色数据',40,7,'','/api/v1/roles/:id/enable','PATCH','01010207','7f6c7556-5242-444f-9714-59a1b5d1abcf',1,1,'root',1543932342,0,0),
	(40,'0568b1ef-8049-4aec-af4b-2041270bcb43','disable','禁用角色数据',40,8,'','/api/v1/roles/:id/disable','PATCH','01010208','7f6c7556-5242-444f-9714-59a1b5d1abcf',1,1,'root',1543932362,0,0),
	(41,'ca78c7b6-8df7-4d09-ba79-a5a33be34451','query','查询用户数据',40,1,'','/api/v1/users','GET','01010301','4b3448fd-c23f-49df-a51b-94aa2a68aec6',1,1,'root',1543933333,0,0),
	(42,'d239aafd-3e04-410e-b3ea-db0fd942a352','one','查询指定用户数据',40,2,'','/api/v1/users/:id','GET','01010302','4b3448fd-c23f-49df-a51b-94aa2a68aec6',1,1,'root',1543933351,0,0),
	(43,'68081a47-53dc-4bde-97e0-b48f83fe043e','create','创建用户数据',40,3,'','/api/v1/users','POST','01010303','4b3448fd-c23f-49df-a51b-94aa2a68aec6',1,1,'root',1543933368,0,0),
	(44,'df0abcc9-d631-4738-aff9-d78da3fb757e','update','更新用户数据',40,4,'','/api/v1/users/:id','PUT','01010304','4b3448fd-c23f-49df-a51b-94aa2a68aec6',1,1,'root',1543933385,0,0),
	(45,'a7a384a7-13dd-4690-affb-e3bc68934434','delete','删除用户数据',40,5,'','/api/v1/users/:id','DELETE','01010305','4b3448fd-c23f-49df-a51b-94aa2a68aec6',1,1,'root',1543933404,0,0),
	(46,'13f59032-fb3d-45a4-83a5-ba79ba98904e','deleteMany','删除多条用户数据',40,5,'','/api/v1/users','DELETE','01010306','4b3448fd-c23f-49df-a51b-94aa2a68aec6',1,1,'root',1543933431,0,0),
	(47,'e54a4e88-1bd7-4186-8fea-892376537925','enable','启用用户数据',40,7,'','/api/v1/users/:id/enable','PATCH','01010307','4b3448fd-c23f-49df-a51b-94aa2a68aec6',1,1,'root',1543933451,0,0),
	(48,'c8e76678-68bc-45b6-917b-df192220b9ae','disable','禁用用户数据',40,8,'','/api/v1/users/:id/disable','PATCH','01010308','4b3448fd-c23f-49df-a51b-94aa2a68aec6',1,1,'root',1543933470,0,0);

UNLOCK TABLES;
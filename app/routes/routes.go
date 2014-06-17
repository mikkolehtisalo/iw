// GENERATED CODE - DO NOT EDIT
package routes

import "github.com/revel/revel"


type tWikis struct {}
var Wikis tWikis


func (_ tWikis) Create(
		wiki string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	return revel.MainRouter.Reverse("Wikis.Create", args).Url
}

func (_ tWikis) Read(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Wikis.Read", args).Url
}

func (_ tWikis) Update(
		wiki string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	return revel.MainRouter.Reverse("Wikis.Update", args).Url
}

func (_ tWikis) Delete(
		wiki string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	return revel.MainRouter.Reverse("Wikis.Delete", args).Url
}


type tFavoriteWikis struct {}
var FavoriteWikis tFavoriteWikis


func (_ tFavoriteWikis) Create(
		wiki string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	return revel.MainRouter.Reverse("FavoriteWikis.Create", args).Url
}

func (_ tFavoriteWikis) Read(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("FavoriteWikis.Read", args).Url
}

func (_ tFavoriteWikis) Delete(
		wiki string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	return revel.MainRouter.Reverse("FavoriteWikis.Delete", args).Url
}


type tUserAvatars struct {}
var UserAvatars tUserAvatars


func (_ tUserAvatars) Read(
		user string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "user", user)
	return revel.MainRouter.Reverse("UserAvatars.Read", args).Url
}


type tUserGroupSearch struct {}
var UserGroupSearch tUserGroupSearch


func (_ tUserGroupSearch) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("UserGroupSearch.List", args).Url
}


type tLocks struct {}
var Locks tLocks


func (_ tLocks) Create(
		wiki string,
		target string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "target", target)
	return revel.MainRouter.Reverse("Locks.Create", args).Url
}

func (_ tLocks) Read(
		wiki string,
		target string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "target", target)
	return revel.MainRouter.Reverse("Locks.Read", args).Url
}

func (_ tLocks) Delete(
		wiki string,
		target string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "target", target)
	return revel.MainRouter.Reverse("Locks.Delete", args).Url
}


type tPages struct {}
var Pages tPages


func (_ tPages) Create(
		wiki string,
		page string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("Pages.Create", args).Url
}

func (_ tPages) Read(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Pages.Read", args).Url
}

func (_ tPages) Update(
		wiki string,
		page string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("Pages.Update", args).Url
}

func (_ tPages) Delete(
		wiki string,
		page string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("Pages.Delete", args).Url
}


type tAttachments struct {}
var Attachments tAttachments


func (_ tAttachments) Create(
		wiki string,
		attachment string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "attachment", attachment)
	return revel.MainRouter.Reverse("Attachments.Create", args).Url
}

func (_ tAttachments) Read(
		wiki string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	return revel.MainRouter.Reverse("Attachments.Read", args).Url
}

func (_ tAttachments) Serve(
		wiki string,
		attachment string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "attachment", attachment)
	return revel.MainRouter.Reverse("Attachments.Serve", args).Url
}

func (_ tAttachments) Delete(
		wiki string,
		attachment string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "attachment", attachment)
	return revel.MainRouter.Reverse("Attachments.Delete", args).Url
}


type tApp struct {}
var App tApp


func (_ tApp) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("App.Index", args).Url
}


type tContentFields struct {}
var ContentFields tContentFields


func (_ tContentFields) Read(
		wiki string,
		page string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("ContentFields.Read", args).Url
}

func (_ tContentFields) Update(
		wiki string,
		page string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "wiki", wiki)
	revel.Unbind(args, "page", page)
	return revel.MainRouter.Reverse("ContentFields.Update", args).Url
}


type tActivities struct {}
var Activities tActivities


func (_ tActivities) Read(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Activities.Read", args).Url
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).Url
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).Url
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).Url
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).Url
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).Url
}



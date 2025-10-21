from django.urls import path
from . import views

urlpatterns = [
    path('', views.device_list, name='device_list'),
    path('create/', views.create_device, name='create_device'),
    path('edit/<str:device_name>/', views.edit_device, name='edit_device'),
    path('delete/<str:device_name>/', views.delete_device, name='delete_device'),
    path('toggle/<str:device_name>/', views.toggle_device, name='toggle_device'),
    path('status/', views.get_device_statuses, name='get_device_statuses'),
]
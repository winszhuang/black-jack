import Swal, { SweetAlertResult } from 'sweetalert2'

const Toast = Swal.mixin({
  toast: true,
  position: 'bottom-end',
  showConfirmButton: false,
  timer: 3000,
  didOpen: (toast) => {
    toast.onmouseenter = Swal.stopTimer
    toast.onmouseleave = Swal.resumeTimer
  }
})

export function notify(message: string): Promise<SweetAlertResult<any>> {
  return Toast.fire({
    icon: 'success',
    title: message
  })
}

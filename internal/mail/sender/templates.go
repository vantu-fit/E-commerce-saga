package sender

import "fmt"

func RegisterEmailTemplate(name string , link string) string {
	content := fmt.Sprintf(`
	<div>
      <h2>Chào mừng %s!</h2>
      <p>Cảm ơn bạn đã đăng ký tài khoản !</p>
      <p>Để kích hoạt tài khoản của bạn, vui lòng nhấp vào nút xác nhận dưới đây:</p>
      <a href="%s">Xác nhận tài khoản</a>
      <p>Lưu ý: Liên kết này chỉ có hiệu lực trong 24 giờ.</p>
    </div>
    <div>
      <p>&copy; All rights reserved.</p>
    </div>
	` , name , link)
	return content
}

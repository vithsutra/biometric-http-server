package utils

func GetUserAccessEmailTemplate(userName string, password string) string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to Vithsutra Technologies Biometric Solution</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'Poppins', Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td style="padding: 0;">
                <table role="presentation" style="width: 100%; max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; margin-top: 20px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <!-- Modern Header Design -->
                    <tr>
                        <td style="background: linear-gradient(135deg, #4169E1 0%, #5c88ff 100%); padding: 40px 30px; border-radius: 8px 8px 0 0; text-align: center;">
                            <div style="margin-bottom: 20px;">
                                <!-- Wave Design -->
                                <div style="margin-bottom: 15px;">
                                    <svg width="100" height="24" viewBox="0 0 100 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M0 12C16.6667 12 16.6667 24 33.3333 24C50 24 50 12 66.6667 12C83.3333 12 83.3333 24 100 24V0H0V12Z" fill="rgba(255,255,255,0.1)"/>
                                    </svg>
                                </div>
                                <h1 style="color: #ffffff; font-size: 28px; margin: 0; font-weight: 600; font-family: 'Poppins', Arial, sans-serif; text-transform: uppercase; letter-spacing: 1px;">
                                    Vithsutra Technologies
                                </h1>
                            </div>
                            <!-- Welcome Message with Modern Design -->
                            <div style="background: rgba(255,255,255,0.1); border-radius: 12px; padding: 25px; margin-top: 20px;">
                                <h2 style="color: #ffffff; font-size: 24px; margin: 0 0 10px 0; font-family: 'Poppins', Arial, sans-serif;">
                                    Welcome to Your Biometric Journey!
                                </h2>
                                <p style="color: rgba(255,255,255,0.9); font-size: 16px; margin: 0; font-family: 'Poppins', Arial, sans-serif; line-height: 1.6;">
                                    Your gateway to secure and efficient biometric solutions
                                </p>
                            </div>
                        </td>
                    </tr>

                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px 30px;">
                            <!-- Welcome Card -->
                            <div style="background: #f8f9fa; border-radius: 12px; padding: 30px; margin-bottom: 30px; border-left: 4px solid #4169E1;">
                                <h3 style="color: #4169E1; font-size: 20px; margin: 0 0 15px 0; font-family: 'Poppins', Arial, sans-serif;">
                                    🎉 Getting Started with Your Account
                                </h3>
                                <p style="color: #666666; font-size: 16px; line-height: 24px; margin: 0; font-family: 'Poppins', Arial, sans-serif;">
                                    We're thrilled to have you join us! Here are your secure login credentials to access our advanced biometric software.
                                </p>
                            </div>
                            
                            <!-- Credentials Box with Modern Design -->
                            <div style="background: linear-gradient(to right, #ffffff, #f8f9fa); border-radius: 12px; padding: 25px; margin: 20px 0; box-shadow: 0 2px 8px rgba(65,105,225,0.1); border: 1px solid rgba(65,105,225,0.2);">
                                <h3 style="color: #4169E1; font-size: 18px; margin: 0 0 20px 0; font-family: 'Poppins', Arial, sans-serif; display: flex; align-items: center;">
                                    <span style="background: #4169E1; color: white; padding: 6px; border-radius: 6px; font-size: 14px; margin-right: 10px;">🔐</span>
                                    Your Secure Access Details
                                </h3>
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 12px 15px; background: rgba(65,105,225,0.05); border-radius: 6px;">
                                            <p style="margin: 0; color: #666666; font-family: 'Poppins', Arial, sans-serif;">Username</p>
                                            <p style="margin: 5px 0 0 0; color: #333333; font-weight: 600; font-family: 'Poppins', Arial, sans-serif; font-size: 16px;">` + userName + `</p>
                                        </td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 12px 15px; margin-top: 10px; display: block; background: rgba(65,105,225,0.05); border-radius: 6px;">
                                            <p style="margin: 0; color: #666666; font-family: 'Poppins', Arial, sans-serif;">Password</p>
                                            <p style="margin: 5px 0 0 0; color: #333333; font-weight: 600; font-family: 'Poppins', Arial, sans-serif; font-size: 16px;">` + password + `</p>
                                        </td>
                                    </tr>
                                </table>
                                <p style="color: #ff6b6b; font-size: 14px; margin: 15px 0 0 0; font-family: 'Poppins', Arial, sans-serif; display: flex; align-items: center;">
                                    <span style="color: #ff6b6b; margin-right: 5px;">⚠️</span>
                                    Please change your password upon first login for security
                                </p>
                            </div>

                            <!-- Download Button with Modern Design -->
                            <div style="text-align: center; margin: 35px 0;">
                                <a href="{{ .DownloadLink }}" style="display: inline-block; padding: 16px 35px; background: linear-gradient(135deg, #4169E1 0%, #5c88ff 100%); color: #ffffff; text-decoration: none; border-radius: 50px; font-family: 'Poppins', Arial, sans-serif; font-weight: 500; transition: all 0.3s ease; box-shadow: 0 4px 15px rgba(65,105,225,0.2);">
                                    Download Your Biometric Software
                                    <span style="margin-left: 5px;">→</span>
                                </a>
                            </div>

                            <!-- Quick Start Guide -->
                            <div style="background: #ffffff; border-radius: 12px; padding: 25px; margin-top: 30px; border: 1px solid #e1e1e1;">
                                <h3 style="color: #333333; font-size: 18px; margin: 0 0 20px 0; font-family: 'Poppins', Arial, sans-serif;">
                                    Quick Start Guide
                                </h3>
                                <div style="display: grid; gap: 15px;">
                                    <div style="display: flex; align-items: start; gap: 15px;">
                                        <span style="background: #4169E1; color: white; width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 14px;">1</span>
                                        <p style="margin: 0; color: #666666; font-size: 15px;">Download and install the software using the button above</p>
                                    </div>
                                    <div style="display: flex; align-items: start; gap: 15px;">
                                        <span style="background: #4169E1; color: white; width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 14px;">2</span>
                                        <p style="margin: 0; color: #666666; font-size: 15px;">Launch the application</p>
                                    </div>
                                    <div style="display: flex; align-items: start; gap: 15px;">
                                        <span style="background: #4169E1; color: white; width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 14px;">3</span>
                                        <p style="margin: 0; color: #666666; font-size: 15px;">Log in with your credentials</p>
                                    </div>
                                    <div style="display: flex; align-items: start; gap: 15px;">
                                        <span style="background: #4169E1; color: white; width: 24px; height: 24px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 14px;">4</span>
                                        <p style="margin: 0; color: #666666; font-size: 15px;">Change your password for security</p>
                                    </div>
                                </div>
                            </div>
                        </td>
                    </tr>

                    <!-- Modern Footer -->
                    <tr>
                        <td style="padding: 30px; background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%); border-radius: 0 0 8px 8px; text-align: center;">
                            <p style="color: #666666; font-size: 14px; margin: 0 0 10px 0; font-family: 'Poppins', Arial, sans-serif;">
                                Need help? Our support team is available 24/7
                            </p>
                            <p style="color: #999999; font-size: 14px; margin: 0; font-family: 'Poppins', Arial, sans-serif;">
                                © 2024 Vithsutra Technologies | Secure Biometric Solutions
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`
}

using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;

namespace PalmBook.Server.Controllers
{
    [Authorize]
    [ApiController]
    [Route("api/account")]
    public class AccountController : ControllerBase
    {
        [HttpGet("userinfo")]
        public IActionResult GetUserInfo()
        {
            var user = HttpContext.User;

            if (user.Identity.IsAuthenticated)
            {
                var userInfo = new
                {
                    Name = user.Identity.Name,
                    Email = user.Claims.FirstOrDefault(c => c.Type == "preferred_username")?.Value,
                    Roles = user.Claims.Where(c => c.Type == "roles").Select(c => c.Value).ToList()
                };

                return Ok(userInfo);
            }

            return Unauthorized();
        }

        [HttpGet("logout")]
        public IActionResult Logout()
        {
            return SignOut("Cookies", "AzureAD");
        }
    }

}

# Admins and Multi-Site Support

Artalk supports multi-site management, allowing administrators to create and manage multiple sites and their respective comments and pages from a centralized interface.

## Creating an Administrator

Run the following command and follow the prompts to quickly create an administrator:

```sh
./artalk admin
```

If you are using Docker, you can run the following command:

```sh
docker exec -it artalk artalk admin
```

After creating the administrator account, you can log into the control panel to create additional accounts through the graphical user interface on the user management page, eliminating the need to manually edit configuration files.

## Creating and Managing Sites

You can create and manage multiple sites and switch between them quickly in the "Dashboard" accessible from the sidebar.

## Admin Configuration

You can set up multiple administrator accounts. When the input field matches an administrator's username and email, a password verification prompt will appear. Only administrators can access the "Dashboard" and manage comments from the frontend.

(Optional) Add administrators via the configuration file:

```yaml
admin_users:
  - name: admin
    email: admin@example.com
    password: (bcrypt)$2y$10$ti4vZYIrxVN8rLcYXVgXCO.GJND0dyI49r7IoF3xqIx8bBRmIBZRm
    badge_name: Administrator
    badge_color: '#0083FF'
  - name: admin2
    email: admin2@example.com
    password: (bcrypt)$2y$10$ti4vZYIrxVN8rLcYXVgXCO.GJND0dyI49r7IoF3xqIx8bBRmIBZRm
    badge_name: Junior Admin
    badge_color: '#0083FF'
```

Explanation of each configuration item:

- **name** & **email**: Username and email, **case-insensitive**.
- **password**: User password.

  Supports bcrypt and md5 encryption. For example, you can specify: `"(md5)50c21190c6e4e5418c6a90d2b5031119"`.

  **Using the more secure bcrypt encryption algorithm is recommended**. In a Linux environment, you can use the [htpasswd command](https://httpd.apache.org/docs/2.4/programs/htpasswd.html) to generate the encrypted password:

  ```bash
  unset HISTFILE # Temporarily disable history to prevent the password from appearing in the history
  htpasswd -bnBC 10 "" "your_password" | tr -d ':'
  ```

  Then configure it as: `"(bcrypt)$2y$10$ti4vZYIrxVN8rLcY..."`, starting with `(bcrypt)`.

  Command explanation reference: [Compute bcrypt hash from command line](https://unix.stackexchange.com/questions/307994/compute-bcrypt-hash-from-command-line#answer-419855)

- **badge_name**: The title badge text displayed for the user.
- **badge_color**: The background color of the title badge displayed for the user.

::: tip

You can also configure administrator accounts through environment variables. Refer to: [Environment Variable Configuration](../env.md)

:::

### Controlling Admin Email Notifications

By default, when there are new comments on a page, emails are sent to all administrators. However, you can configure `receive_email` to forcefully disable email notifications for specific administrators.

This is useful if you have multiple email addresses configured but do not want certain addresses to receive comment notifications.

- **receive_email**: When set to `false`, the system will not send email notifications to that user.

  Note: Even if email notifications are disabled, the administrator will still receive @AT replies from others. The system will not send email notifications when the user comments on a page (creates a root comment).

```yaml
admin_users:
  - name: admin
    receive_email: false # ← Forcefully disable email notifications
```

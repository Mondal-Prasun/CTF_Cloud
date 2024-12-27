FROM kalilinux/kali-rolling

# Install VNC server, noVNC, and tools
RUN apt-get update && apt-get install -y \
    tightvncserver \
    novnc \
    xauth \
    xfce4 \
    xfce4-terminal \
    dbus-x11 \
    net-tools \
    curl \
    && apt-get clean

# Set environment variables
ENV USER=root
ENV HOME=/root
ENV DISPLAY=:1

# Create .Xauthority file
RUN touch /root/.Xauthority

# Create necessary directories and files
RUN mkdir -p /root/.vnc && \
    echo "startxfce4" > /root/.vnc/xstartup && \
    chmod +x /root/.vnc/xstartup

# Add a default password for VNC (replace 'password' with your desired password)
RUN mkdir -p /root/.vnc && echo "123" | vncpasswd -f > /root/.vnc/passwd && chmod 600 /root/.vnc/passwd

# Expose the noVNC port
EXPOSE 8080

# Remove CMD to manually run the container
